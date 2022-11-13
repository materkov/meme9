import React from "react";
import {useQuery} from "@tanstack/react-query";
import {Photo} from "../store/types";
import {fetcher} from "../store/fetcher";

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

    const maxWidth = 400;
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

    const styles = {
        width: data.width * ratio + 'px',
        height: data.height * ratio + 'px',
    };

    return <img alt={""} src={data.address} onClick={onPhotoClick} style={styles} className={className}/>
}
