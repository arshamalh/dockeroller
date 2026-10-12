package handlers

import (
	"strings"

	"github.com/arshamalh/dockeroller/entities"
	"github.com/arshamalh/dockeroller/telegram/keyboards"
	"github.com/arshamalh/dockeroller/telegram/msgs"
	"github.com/arshamalh/dockeroller/tools"
)

func (h *Handler) refreshImages(c Ctx) ([]*entities.Image, error) {
	ctx, cancel := h.dockerCtx()
	defer cancel()
	images, err := h.docker.ImagesList(ctx)
	if err != nil {
		return nil, err
	}
	c.State().Images = images
	return images, nil
}

func showImage(c Ctx, index int, img *entities.Image, prefix string) error {
	text := c.GetL8d(msgs.Image, msgs.ImageArgs(img)...)
	if prefix != "" {
		text = prefix + "\n\n" + text
	}
	kb := keyboards.Image(c, index)
	if c.Callback() != nil {
		return c.CleanEdit(text, kb)
	}
	return c.CleanSend(text, kb)
}

func noImages(c Ctx) error {
	if c.Callback() != nil {
		return c.RespondL8d(msgs.ToastNoImages)
	}
	return c.SendL8d(msgs.NoImages)
}

// ImagesList shows the first image (command or menu button).
func (h *Handler) ImagesList(c Ctx) error {
	h.stopStream(c.ID())
	c.ResetScene()
	images, err := h.refreshImages(c)
	if err != nil {
		return dockerUnavailable(c, err)
	}
	if len(images) == 0 {
		return noImages(c)
	}
	if c.Callback() != nil {
		_ = c.Respond()
	}
	return showImage(c, 0, images[0], "")
}

// ImagesNav handles Prev and Next.
func (h *Handler) ImagesNav(c Ctx) error {
	index, err := buttonIndex(c)
	if err != nil {
		return c.RespondL8d(msgs.ToastInvalidButton)
	}
	images, err := h.refreshImages(c)
	if err != nil {
		return dockerUnavailable(c, err)
	}
	if len(images) == 0 {
		return noImages(c)
	}
	index = tools.Indexer(index, len(images))
	_ = c.Respond()
	return showImage(c, index, images[index], "")
}

// cachedImage is cachedContainer for images.
func (h *Handler) cachedImage(c Ctx) (index int, img *entities.Image, ok bool, err error) {
	index, err = buttonIndex(c)
	if err != nil {
		return 0, nil, false, c.RespondL8d(msgs.ToastInvalidButton)
	}
	images := c.State().Images
	if index >= 0 && index < len(images) {
		return index, images[index], true, nil
	}
	fresh, err := h.refreshImages(c)
	if err != nil {
		return 0, nil, false, dockerUnavailable(c, err)
	}
	if len(fresh) == 0 {
		return 0, nil, false, noImages(c)
	}
	_ = c.RespondL8d(msgs.ToastListChanged)
	return 0, nil, false, showImage(c, 0, fresh[0], "")
}

// ImageBack returns to the image view and leaves the tag scene.
func (h *Handler) ImageBack(c Ctx) error {
	c.ResetScene()
	index, img, ok, err := h.cachedImage(c)
	if !ok {
		return err
	}
	_ = c.Respond()
	return showImage(c, index, img, "")
}

// ImageTag asks for a new tag; the next text message answers it.
func (h *Handler) ImageTag(c Ctx) error {
	index, img, ok, err := h.cachedImage(c)
	if !ok {
		return err
	}
	c.Enter(SceneTagImage, QNewImageTag)
	c.State().CurrentImage = img
	_ = c.Respond()
	return c.CleanEdit(c.GetL8d(msgs.ImageNewTag, "id", msgs.ShortImageID(img.ID)), keyboards.ImageBack(c, index))
}

// ImageTagText tags the image with the text the user sent. On failure the
// scene stays so another tag can be sent; on success it ends.
func (h *Handler) ImageTagText(c Ctx) error {
	img := c.State().CurrentImage
	if img == nil {
		c.ResetScene()
		return h.Start(c)
	}
	tag := strings.TrimSpace(c.Text())

	ctx, cancel := h.dockerCtx()
	defer cancel()
	if err := h.docker.ImageTag(ctx, img.ID, tag); err != nil {
		c.Logger().Warn("tag failed", "image", img.ID, "err", err)
		return c.CleanSend(
			c.GetL8d(msgs.ImageTagErr, "id", msgs.ShortImageID(img.ID), "error", err.Error()),
			keyboards.ImageBack(c, h.imageIndex(c, img.ID)),
		)
	}

	c.ResetScene()
	c.State().CurrentImage = nil
	images, err := h.refreshImages(c)
	if err != nil {
		return dockerUnavailable(c, err)
	}
	if len(images) == 0 {
		return c.SendL8d(msgs.NoImages)
	}
	index := 0
	for i, x := range images {
		if x.ID == img.ID {
			index = i
			break
		}
	}
	done := c.GetL8d(msgs.ImageTagged, "id", msgs.ShortImageID(img.ID), "tag", tag)
	return showImage(c, index, images[index], done)
}

func (h *Handler) imageIndex(c Ctx, id string) int {
	for i, x := range c.State().Images {
		if x.ID == id {
			return i
		}
	}
	return 0
}

// ImageRemoveForm shows the remove options (reset for every image) and
// waits for Done.
func (h *Handler) ImageRemoveForm(c Ctx) error {
	index, img, ok, err := h.cachedImage(c)
	if !ok {
		return err
	}
	c.State().ImageRemove = entities.ImageRemoveForm{}
	_ = c.RespondL8d(msgs.ToastFillFormPressDone)
	return h.showImageRemoveForm(c, index, img)
}

func (h *Handler) showImageRemoveForm(c Ctx, index int, img *entities.Image) error {
	text := c.GetL8d(msgs.RemoveFormHint) + "\n\n" + c.GetL8d(msgs.Image, msgs.ImageArgs(img)...)
	return c.CleanEdit(text, keyboards.ImageRemove(c, index, c.State().ImageRemove))
}

func (h *Handler) ImageRemoveForce(c Ctx) error {
	index, img, ok, err := h.cachedImage(c)
	if !ok {
		return err
	}
	form := &c.State().ImageRemove
	form.Force = !form.Force
	_ = c.Respond()
	return h.showImageRemoveForm(c, index, img)
}

func (h *Handler) ImageRemovePruneChildren(c Ctx) error {
	index, img, ok, err := h.cachedImage(c)
	if !ok {
		return err
	}
	form := &c.State().ImageRemove
	form.PruneChildren = !form.PruneChildren
	_ = c.Respond()
	return h.showImageRemoveForm(c, index, img)
}

// ImageRemoveDone removes the image with the chosen options.
func (h *Handler) ImageRemoveDone(c Ctx) error {
	_, img, ok, err := h.cachedImage(c)
	if !ok {
		return err
	}
	form := c.State().ImageRemove
	ctx, cancel := h.dockerCtx()
	defer cancel()
	if err := h.docker.ImageRemove(ctx, img.ID, form.Force, form.PruneChildren); err != nil {
		c.Logger().Warn("image remove failed", "image", img.ID, "err", err)
		return respondErr(c, err)
	}
	c.State().ImageRemove = entities.ImageRemoveForm{}
	_ = c.RespondL8d(msgs.ToastImageRemoved)

	images, err := h.refreshImages(c)
	if err != nil {
		return dockerUnavailable(c, err)
	}
	if len(images) == 0 {
		return c.CleanEdit(c.GetL8d(msgs.NoImages), keyboards.Welcome(c))
	}
	return showImage(c, 0, images[0], "")
}
