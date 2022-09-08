import React, {useEffect} from "react";
import styles from "./Header.module.css";
import {Link} from "./Link";
import {User} from "../store/types";

function getViewer(): Promise<User> {
    return new Promise<User>((resolve, reject) => {

        fetch("http://localhost:8000/browse?url=/viewer", {
            headers: {
                'authorization': 'Bearer ' + localStorage.getItem('authToken'),
            }
        })
            .then(r => r.json())
            .then(r => {
                resolve(r.componentData[0] as User)
            })
    })

}

const vkURL = "https://oauth.vk.com/authorize?client_id=7260220&response_type=code&redirect_uri=http://localhost:3000/vk-callback"

export function Header() {
    const [viewer, setViewer] = React.useState<User | undefined | null>();
    useEffect(() => {
        getViewer().then(setViewer);
    }, [])

    const onLogout = (e: MouseEvent) => {
        e.stopPropagation();
        localStorage.removeItem('authToken');

        setViewer(null);
    }

    return (
        <div className={styles.header}>
            <Link href="/" className={styles.logo}>meme</Link>


            <div className={styles.userName}>
                {viewer === null && <a href={vkURL}>Авторизация</a>}
                {viewer !== null && viewer !== undefined &&
                    <>
                    <Link href={viewer.href}>{viewer.name}</Link> | <a onClick={onLogout} href={"#"}>Выход</a>
                    </>
                }
            </div>
        </div>
    )
}