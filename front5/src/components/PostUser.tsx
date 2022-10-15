import React from "react";
import {api, Post, User} from "../store/types";
import {Link} from "./Link";
import styles from "./PostUser.module.css";
import {UserAvatar} from "./UserAvatar";

export function PostUser(props: { post: Post }) {
    const [isVisible, setIsVisible] = React.useState(false);
    const [userName, setUserName] = React.useState("");
    const [postsCount, setPostsCount] = React.useState(0);
    const [isUserLoaded, setIsUserLoaded] = React.useState(false);

    let className = styles.userNamePopup;
    if (!isVisible) {
        className += " " + styles.userNamePopup__hidden;
    }

    const onMouseEnter = () => {
        setIsVisible(true);

        if (isUserLoaded) return;

        setIsUserLoaded(true);

        const f = new FormData();
        f.set("id", props.post.user?.id || "");

        api("/users.list", {
            ids: props.post.user?.id || "",
            fields: "posts",
        }).then(r => {
            setUserName(r[0].name);
            setPostsCount(r[0].posts.count);
        })
    }

    let userDetails = '...LOADING...';
    if (userName) {
        userDetails = "Name: " + userName + ", posts: " + postsCount;
    }

    const date = new Date(props.post.date || "");
    const dateStr = date.toLocaleString();

    return (
        <div className={styles.userName}
             onMouseEnter={onMouseEnter}
             onMouseLeave={() => setIsVisible(false)}
        >
            <div className={className}>{userDetails}</div>

            <UserAvatar width={50} url={props.post.user?.avatar}/>
            <div className={styles.rightContainer}>
                <Link href={"/users/" + props.post.user?.id} className={styles.name}>{props.post.user?.name}</Link>
                <Link href={"/posts/" + props.post.id} className={styles.href}>{dateStr}</Link>
            </div>
        </div>
    )
}
