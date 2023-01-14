import React, {MouseEvent, useEffect} from "react";
import styles from "./Header.module.css";
import {Link} from "./Link";
import {authorize} from "../utils/localize";
import {Global} from "../store2/store";
import {actions} from "../store2/actions";
import {connect} from "react-redux";
import * as types from "../store/types";

interface Props {
    isLoaded: boolean;
    viewerId: string;
    viewer?: types.User;
}

function Component(props: Props) {
    const [data, setData] = React.useState(false);

    useEffect(() => {
        actions.loadViewer().then(() => setData(true));
    }, []);

    const onLogout = (e: MouseEvent<HTMLAnchorElement>) => {
        e.preventDefault();

        authorize('');
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