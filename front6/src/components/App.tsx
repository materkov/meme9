import React, {useEffect} from "react";
import {Profile} from "./Profile/Profile";
import {ReactQueryDevtools} from '@tanstack/react-query-devtools'
import * as styles from "./App.module.css";
import {useNavigation} from "../store/navigation";
import {Link} from "./Link/Link";
import {useGlobals} from "../store/globals";
import {Auth} from "./Auth/Auth";
import {Discover} from "./Discover/Discover";
import {QueryClientProvider} from "@tanstack/react-query";
import {PostPage} from "./PostPage/PostPage";
import {BookmarksPage} from "./BookmarksPage/BookmarksPage";
import {queryClient} from "../store/queryClient";


export function App() {
    useEffect(() => {
        const serverRender = document.querySelector('#server-render');
        if (serverRender) {
            serverRender.parentElement?.removeChild(serverRender);
        }
    }, []);

    return (
        <QueryClientProvider client={queryClient}>
            <AppContent/>
            <ReactQueryDevtools initialIsOpen={true}/>
        </QueryClientProvider>
    )
}

export function AppContent() {
    let page: React.ReactNode;

    const globals = useGlobals();

    const navState = useNavigation();

    // Kinda router
    if (navState.url === "/") {
        page = <Discover/>
    } else if (navState.url.startsWith("/posts/")) {
        page = <PostPage/>
    } else if (navState.url.startsWith("/users/")) {
        page = <Profile/>
    } else if (navState.url.startsWith("/auth")) {
        page = <Auth/>
    } else if (navState.url.startsWith("/bookmarks")) {
        page = <BookmarksPage/>
    } else if (navState.url.startsWith("/vk-callback")) {
        page = <Auth/>
    } else {
        page = <div>404 page</div>;
    }

    return (
        <div className={styles.app}>
            <div className={styles.header}>
                <Link href={"/"} className={styles.headerLink}>
                    meme
                </Link>

                <div className={styles.authInfo}>
                    {!globals.viewerId && <Link href="/auth">Authorize</Link>}

                    {globals.viewerId &&
                        <span>
                                <Link href={"/users/" + globals.viewerId}>{globals.viewerName}</Link>
                            &nbsp;|&nbsp;
                            <Link href="/auth?logout">Logout</Link>
                            </span>
                    }
                </div>
            </div>

            <div className={styles.nav}>
                <Link href={"/"}>Index page</Link>
                &nbsp;|&nbsp;
                <Link href={"/bookmarks"}>Bookmarks</Link>
            </div>

            {page}
        </div>
    )
}