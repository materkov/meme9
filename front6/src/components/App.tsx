import React, {useEffect} from "react";
import {ArticlePage} from "./ArticlePage";
import {Profile} from "./Profile/Profile";
import * as styles from "./App.module.css";
import {useNavigation} from "../store/navigation";
import {Discover} from "./Discover/Discover";
import {Link} from "./Link";
import {useGlobals} from "../store/globals";

export function App() {
    let page: React.ReactNode;

    const globals = useGlobals();

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

    useEffect(() => {
        const serverRender = document.querySelector('#server-render');
        if (serverRender) {
            serverRender.parentElement?.removeChild(serverRender);
        }
    }, []);

    const redirectURL = location.origin + "/vk-callback";
    const vkAuthURL = "https://oauth.vk.com/authorize?client_id=7260220&response_type=code&v=5.131&redirect_uri=" + redirectURL;

    const onLogout = (e: React.MouseEvent<HTMLAnchorElement>) => {
        globals.logout();
        e.preventDefault();
    };

    return <div className={styles.app}>
        <div className={styles.header}>
            <Link href={"/"} className={styles.headerLink}>
                meme
            </Link>

            <div className={styles.authInfo}>
                {!globals.viewerId && <a href={vkAuthURL}>Войти через VK</a>}
                {globals.viewerId && <span>
                    {globals.viewerName} | <a href="/logout" onClick={onLogout}>Выйти</a>
                </span>}
            </div>
        </div>

        {page}
    </div>;
}
