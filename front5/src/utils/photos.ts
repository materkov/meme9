import {Photo} from "../store/types";

export function selectPhotoThumb(photo: Photo, width: number): string {
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
