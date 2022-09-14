import React from "react";
import {api, Post, User} from "../store/types";
import {Link} from "./Link";
import styles from "./PostUser.module.css";

const avatarStub = 'https://sun9-73.userapi.com/s/v1/ig2/dv3b5tV5Umau1mxqiWJwb6hJOHc-f5_lEkNmmjWuZ3_hsVQcfH9yiril0lJTbKjDr1Hc9BBZU_RY6aldGSU8N9cR.jpg?size=100x100&quality=95&crop=340,512,1228,1228&ava=1';

export function PostUser(props: { post: Post }) {
    const [isVisible, setIsVisible] = React.useState(false);
    const [userData, setUserData] = React.useState<User>();
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

        api("/userPage", {
            id: props.post.user?.id || ""
        }).then(r => setUserData(r[0]));
    }

    let userDetails = '...LOADING...';
    if (userData) {
        userDetails = "Name: " + userData.name + ", posts: " + userData.posts?.count;
    }

    const date = new Date(props.post.date || "");
    const dateStr = date.toLocaleString();

    return (
        <div className={styles.userName}
             onMouseEnter={onMouseEnter}
             onMouseLeave={() => setIsVisible(false)}
        >
            <div className={className}>{userDetails}</div>

            <img src={props.post.user?.avatar} className={styles.avatar} alt=""/>
            <div className={styles.rightContainer}>
                <Link href={"/users/" + props.post.user?.id} className={styles.name}>{props.post.user?.name}</Link>
                <Link href={"/posts/" + props.post.id} className={styles.href}>{dateStr}</Link>
            </div>
        </div>
    )
}
