import React from "react";
import {User} from "../store/types";
import {Link} from "./Link";

export function PostUser(props: { user: User }) {
    return (
        <div>
            From: <Link href={props.user.href}>{props.user.name}</Link>
        </div>
    )
}
