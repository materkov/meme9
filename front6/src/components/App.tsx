import React from "react";
import {ArticlePage} from "./ArticlePage";
import {Profile} from "./Profile/Profile";
import * as styles from "./App.module.css";
import {useNavigation} from "../store/navigation";
import {Discover} from "./Discover/Discover";
import {Link} from "./Link";

export function App() {
    let page: React.ReactNode;

    const navState = useNavigation(state => state);

    if (navState.url === "/") {
        page = <Discover/>
    } else if (navState.url.startsWith("/article/")) {
        page = <ArticlePage/>
    } else if (navState.url.startsWith("/users/")) {
        page = <Profile/>
    } else {
        page = <div>404 page</div>;
    }

    return <div className={styles.app}>
        <div className={styles.header}>
            <Link href={"/"} className={styles.headerLink}>
                meme
            </Link>
        </div>

        {page}
    </div>;
}
