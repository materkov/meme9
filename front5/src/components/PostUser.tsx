import React from "react";
import {BrowseResult, Post, User} from "../store/types";
import {Link} from "./Link";
import styles from "./PostUser.module.css";

export function PostUser(props: { user: User }) {
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
        fetch("http://localhost:8000/browse?url=/users/" + props.user.id)
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

    return (
        <div className={styles.userName}
             onMouseEnter={onMouseEnter}
             onMouseLeave={() => setIsVisible(false)}
        >
            <div className={className}>{userDetails}</div>
            From: <Link href={props.user.href}>{props.user.name}</Link>
        </div>
    )
}
