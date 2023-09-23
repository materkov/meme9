import React, {useEffect} from "react";
import {PostPage} from "./PostPage/PostPage";
import {Profile} from "./Profile/Profile";
import * as styles from "./App.module.css";
import {useNavigation} from "../store/navigation";
import {Discover} from "./Discover/Discover";
import {Link} from "./Link/Link";
import {useGlobals} from "../store/globals";
import {Auth} from "./Auth/Auth";

export function App() {
    let page: React.ReactNode;

    const globals = useGlobals();

    const navState = useNavigation(state => state);

    // Kinda router
    if (navState.url === "/") {
        page = <Discover/>
    } else if (navState.url.startsWith("/posts/")) {
        page = <PostPage/>
    } else if (navState.url.startsWith("/users/")) {
        page = <Profile/>
    } else if (navState.url.startsWith("/auth")) {
        page = <Auth/>
    } else if (navState.url.startsWith("/vk-callback")) {
        page = <Auth/>
    } else {
        page = <div>404 page</div>;
    }

    useEffect(() => {
        const serverRender = document.querySelector('#server-render');
        if (serverRender) {
            serverRender.parentElement?.removeChild(serverRender);
        }
    }, []);

    return <div className={styles.app}>
        <div className={styles.header}>
            <Link href={"/"} className={styles.headerLink}>
                meme
            </Link>

            <div className={styles.authInfo}>
                {!globals.viewerId && <Link href="/auth">Authorize</Link>}
                {globals.viewerId && <span>
                    <Link href={"/users/" + globals.viewerId}>{globals.viewerName}</Link>
                    &nbsp;|&nbsp;
                    <Link href="/auth?logout">Logout</Link>
                </span>}
            </div>
        </div>

        {page}
    </div>;
}
