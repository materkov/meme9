import React from "react";
import {useQuery} from "@tanstack/react-query";
import {Photo} from "../store/types";
import {fetcher} from "../store/fetcher";
import {selectPhotoThumb} from "../utils/photos";
import classNames from "classnames";
import styles from "./PostPhoto.module.css";

export interface Props {
    className: string;
    id: string;
}

export const PostPhoto: React.FC<Props> = ({id, className}) => {
    const {data} = useQuery<Photo>(["/photos/" + id], fetcher);

    const onPhotoClick = () => {
        if (data && data.address) {
            window.open(data.address, '_blank');
        }
    }

    if (!data) return null;

    const maxWidth = 100;
    const maxHeight = 300;

    let ratioWidth = maxWidth / data.width;
    let rationHeight = maxHeight / data.height;

    if (ratioWidth > 1) {
        ratioWidth = 1;
    }
    if (rationHeight > 1) {
        rationHeight = 1;
    }

    const ratio = Math.min(ratioWidth, rationHeight);

    const inlineStyles = {
        width: data.width * ratio + 'px',
        height: data.height * ratio + 'px',
    };

    const url = selectPhotoThumb(data, maxWidth);

    return <img alt={""} src={url} onClick={onPhotoClick} style={inlineStyles}
                className={classNames([className, styles.img])}
    />
}
