import React from "react";
import {apiHost, BrowseResult, Post, User} from "../store/types";
import {Link} from "./Link";
import styles from "./PostUser.module.css";

const avatarStub = 'https://sun9-73.userapi.com/s/v1/ig2/dv3b5tV5Umau1mxqiWJwb6hJOHc-f5_lEkNmmjWuZ3_hsVQcfH9yiril0lJTbKjDr1Hc9BBZU_RY6aldGSU8N9cR.jpg?size=100x100&quality=95&crop=340,512,1228,1228&ava=1';

export function PostUser(props: { post: Post }) {
    const [isVisible, setIsVisible] = React.useState(false);
    const [userData, setUserData] = React.useState<BrowseResult>();
    const [isUserLoaded, setIsUserLoaded] = React.useState(false);

    let className = styles.userNamePopup;
    if (!isVisible) {
        className += " " + styles.userNamePopup__hidden;
    }

    const onMouseEnter = () => {
        setIsVisible(true);

        if (isUserLoaded) return;

        setIsUserLoaded(true);
        fetch(apiHost + "/browse?url=/users/" + props.post.from?.id)
            .then(r => r.json())
            .then(r => setUserData(r))
    }

    let userDetails = '...LOADING...';
    if (userData && userData.componentData) {
        const userObject = userData.componentData[0] as User;
        const userPosts = userData.componentData[1] as Post[];
        if (userObject) {
            userDetails = "Name: " + userObject.name + ", posts: " + userPosts?.length;
        }
    }

    const date = new Date(props.post.date || "");
    const dateStr = date.toLocaleString();

    return (
        <div className={styles.userName}
             onMouseEnter={onMouseEnter}
             onMouseLeave={() => setIsVisible(false)}
        >
            <div className={className}>{userDetails}</div>

            <img src={avatarStub} className={styles.avatar}/>
            <div className={styles.rightContainer}>
                <Link href={props.post.from?.href} className={styles.name}>{props.post.from?.name}</Link>
                <Link href={props.post.detailsURL} className={styles.href}>{dateStr}</Link>
            </div>
        </div>
    )
}
