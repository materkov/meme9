import React from "react";
import * as schema from "../../api/login";

export function Login(props: { data: schema.LoginPageRenderer }) {
    return (
        <>
            <a href={props.data.authUrl}>{props.data.text}</a>
        </>
    )
}
