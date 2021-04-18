import React from "react";
import styles from "./Header.module.css";
import {Link} from "../Link/Link";
import * as schema from "../../api/api2";

export class Header extends React.Component<{ data: schema.HeaderRenderer }> {
    render() {
        const data = this.props.data;

        return (
            <div className={styles.Header}>
                <Link href={data.mainUrl} className={styles.Logo}>meme</Link>

                <div className={styles.RightContainer}>
                    <div className={styles.Name}>{data.userName}</div>
                    <img className={styles.Avatar} alt="" src={data.userAvatar}/>
                </div>
            </div>
        );
    }
}
