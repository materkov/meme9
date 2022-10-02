import React, {MouseEvent, useEffect} from "react";
import styles from "./Header.module.css";
import {Link} from "./Link";
import {api, User} from "../store/types";
import {useCustomEventListener} from "react-custom-events";
import {authorize} from "../utils/localize";

export function Header() {
    const [viewer, setViewer] = React.useState<User | undefined | null>();

    const onLogout = (e: MouseEvent<HTMLAnchorElement>) => {
        e.preventDefault();

        authorize('');
        setViewer(null);
    }

    const refreshUser = () => {
        api("/viewer").then(r => {
            setViewer(r[0]);
        })
    }

    useCustomEventListener('onAuthorized', refreshUser);
    useEffect(refreshUser, [])

    return (
        <div className={styles.header}>
            <Link href="/" className={styles.logo}>meme</Link>

            <div className={styles.userName}>
                {viewer === null && <Link href={"/login"}>Авторизация</Link>}
                {viewer !== null && viewer !== undefined &&
                    <>
                        <Link href={"/users/" + viewer.id}>{viewer.name}</Link> | <a onClick={onLogout} href={"#"}>Выход</a>
                    </>
                }
            </div>
        </div>
    )
}