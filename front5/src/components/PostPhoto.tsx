import React from "react";
import {selectPhotoThumb} from "../utils/photos";
import classNames from "classnames";
import styles from "./PostPhoto.module.css";
import {Global} from "../store/store";
import * as types from "../api/types";
import {connect} from "react-redux";

export interface OwnProps {
    className: string;
    id: string;
}

interface Props {
    className: string;
    photo: types.Photo;
}

export const Component: React.FC<Props> = (props: Props) => {
    const onPhotoClick = () => {
        if (props.photo.address) {
            window.open(props.photo.address, '_blank');
        }
    }

    const maxWidth = 100;
    const maxHeight = 300;

    let ratioWidth = maxWidth / props.photo.width;
    let rationHeight = maxHeight / props.photo.height;

    if (ratioWidth > 1) {
        ratioWidth = 1;
    }
    if (rationHeight > 1) {
        rationHeight = 1;
    }

    const ratio = Math.min(ratioWidth, rationHeight);

    const inlineStyles = {
        width: props.photo.width * ratio + 'px',
        height: props.photo.height * ratio + 'px',
    };

    const url = selectPhotoThumb(props.photo, maxWidth);

    return <img alt={""} src={url} onClick={onPhotoClick} style={inlineStyles}
                className={classNames([props.className, styles.img])}
    />
}

export const PostPhoto = connect((state: Global, ownProps: OwnProps): Props => {
    return {
        className: ownProps.className,
        photo: state.photos.byId[ownProps.id],
    }
})(Component);