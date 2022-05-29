import React from "react";
import {HeaderRenderer} from "./types";
import {Link} from "./Link";

export const Header = (props: { data: HeaderRenderer }) => {
    return (
        <>
            <h1>meme</h1>

            {props.data.userName && <p>Авторизован как {props.data.userName}</p>}
            {!props.data.userName && <Link href={"/login"}>Войти</Link>}
            <hr/>
        </>
    );
}
