import React from "react";
import styles from "./Header.module.css";
import {Link} from "../Link/Link";
import {HeaderRenderer} from "../../api/api2";

export function Header(props: { data: HeaderRenderer }) {
    const data = props.data;

    return (
        <div className={styles.Header}>
            <Link href={data.mainUrl} className={styles.Logo}>meme</Link>

            <div className={styles.RightContainer}>
                {data.isAuthorized && <div className={styles.Name}>{data.userName}</div>}
                {data.isAuthorized && <img className={styles.Avatar} alt="" src={data.userAvatar}/>}
                {data.isAuthorized &&
                <a href={data.logoutUrl}>Выход</a>
                }

                {!data.isAuthorized && <Link className={styles.Name} href={"/login"}>Войти</Link>}
            </div>
        </div>
    );
}
