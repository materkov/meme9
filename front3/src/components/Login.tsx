import React from "react";
import {LoginRenderer} from "./types";
import {Link} from "./Link";

export const Login = (props: { data: LoginRenderer }) => {
    return <>
        <Link href={props.data.authURL}>Войти</Link>
    </>
}
