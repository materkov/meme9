import React, {MouseEvent} from "react";
import styles from "./Header.module.css";
import {Link} from "./Link";
import {User, Viewer} from "../store/types";
import {authorize} from "../utils/localize";
import {fetcher, queryClient} from "../store/fetcher";
import { useQuery } from "@tanstack/react-query";

export function Header() {
    const {data: viewer, isLoading} = useQuery<Viewer>(["/viewer"], fetcher);
    const {data: viewerUser} = useQuery<User>(["/users/" + viewer?.viewerId], fetcher, {
        enabled: !!viewer?.viewerId,
    })
    //const [viewer, setViewer] = React.useState<User | undefined | null>();

    const onLogout = (e: MouseEvent<HTMLAnchorElement>) => {
        e.preventDefault();

        authorize('');
        queryClient.invalidateQueries(["/viewer"]);
    }

    //useCustomEventListener('onAuthorized', refreshUser);
    //useEffect(refreshUser, [])

    return (
        <div className={styles.header}>
            <Link href="/" className={styles.logo}>meme</Link>

            <div className={styles.userName}>
                {!isLoading && !viewer?.viewerId && <Link href={"/login"}>Авторизация</Link>}
                {!isLoading && viewer?.viewerId &&
                    <>
                        <Link href={"/users/" + viewer.viewerId}>{viewerUser?.name}</Link> | <a onClick={onLogout}
                                                                                     href={"#"}>Выход</a>
                    </>
                }
            </div>
        </div>
    )
}