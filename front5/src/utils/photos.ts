import * as types from "../api/types";

export function selectPhotoThumb(photo: types.Photo, width: number): string {
    const actualWidth = Math.floor(width * window.devicePixelRatio);

    if (photo.thumbs) {
        for (let thumb of photo.thumbs) {
            if (thumb.width >= actualWidth) {
                return thumb.address;
            }
        }
    }

    return photo.address;
}
