import * as types from "../../api/api";
import React from "react";

export function PhotoAttach(props: {photo: types.File}) {
    let height = props.photo.height;
    let width = props.photo.width;
    let ratio = 0;

    if (height > 200) {
        ratio = 200 / height;
        height = 200;
        width = Math.floor(width * ratio);
    }

    const styles = {
        width: width + 'px',
        height: height + 'px',
    }

    return <img src={props.photo.url} style={styles}/>
}