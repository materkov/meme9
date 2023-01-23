import React, {MouseEvent, useEffect} from "react";
import styles from "./Header.module.css";
import {Link} from "./Link";
import {Global} from "../store/store";
import {connect} from "react-redux";
import * as types from "../api/types";
import {loadViewer, logout} from "../store/actions/auth";

interface Props {
    isLoaded: boolean;
    viewerId: string;
    viewer?: types.User;
}

function Component(props: Props) {
    useEffect(() => {
        loadViewer();
    }, []);

    const onLogout = (e: MouseEvent<HTMLAnchorElement>) => {
        e.preventDefault();

        logout();
    }

    return (
        <div className={styles.header}>
            <Link href="/" className={styles.logo}>meme</Link>

            <div className={styles.userName}>
                {props.isLoaded && !props.viewerId && <Link href={"/login"}>Авторизация</Link>}
                {props.isLoaded && props.viewerId &&
                    <>
                        <Link
                            href={"/users/" + props.viewerId}>{props.viewer?.name}</Link> | <a
                        onClick={onLogout} href={"#"}>Выход</a>
                    </>
                }
            </div>
        </div>
    )
}

export const Header = connect((state: Global) => {
    return {
        isLoaded: state.viewer.isLoaded,
        viewerId: state.viewer.id,
        viewer: state.viewer.id ? state.users.byId[state.viewer.id] : undefined,
    } as Props;
})(Component);