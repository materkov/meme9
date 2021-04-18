import React from "react";
import {Link} from "../Link/Link";
import * as schema from "../../api/api2";

export function Profile(props: { data: schema.ProfileRenderer }) {
    return <div>
        Profile page
        <br/>
        <Link href={"/"}>Go to feed</Link>
        <br/>
        <img alt="" src={props.data.avatar}/>
        <br/>
        <h2>{props.data.name}</h2>
    </div>;
}
