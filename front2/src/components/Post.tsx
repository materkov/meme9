import React from "react";
import {PostRenderer} from "../types";
import {Link} from "./Link";

export const Post = (props: { data: PostRenderer }) => {
    return (
        <div>
            From: <b><Link href={props.data.authorHref}>{props.data.authorName}</Link></b>
            <br/>
            {props.data.text}
        </div>
    );
}
