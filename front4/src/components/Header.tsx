import styles from "./Header.module.css";
import React, {useEffect} from "react";
import {QueryParams, User} from "../types";
import {api} from "../api";

const localStorageKey = 'authUser';

export function Header() {
    const [authURL, setAuthURL] = React.useState('');
    const [isLoaded, setIsLoaded] = React.useState(false);
    const [viewer, setViewer] = React.useState<User | undefined>();

    useEffect(() => {
        const authUser = localStorage.getItem(localStorageKey);
        if (authUser) {
            setViewer(JSON.parse(authUser));
            setIsLoaded(true);
        }

        const q: QueryParams = {
            viewer: {
                inner: {
                    name: {}
                }
            }
        }
        api(q).then(result => {
            if (result.viewer) {
                setViewer(result.viewer);
                localStorage.setItem(localStorageKey, JSON.stringify(result.viewer));
            }
            setIsLoaded(true);
        })

        const urlQuery: QueryParams = {
            vkAuthUrl: {}
        };
        api(urlQuery).then(result => {
            setAuthURL(result.vkAuthUrl || "");
        })
    }, []);

    const onExit = (e: React.MouseEvent<HTMLAnchorElement>) => {
        e.preventDefault();

        localStorage.removeItem(localStorageKey);
        localStorage.removeItem('authToken');
        setViewer(undefined);
    };

    return (
        <div className={styles.header}>
            <a href={"/"} className={styles.headerLink}>meme9</a>

            {isLoaded &&
                <>
                    {viewer && <div>{viewer.name} | <a href={"#"} onClick={onExit}>Выйти</a></div>}
                    {!viewer && <a href={authURL}>Войти</a>}
                </>
            }
        </div>
    )
}
