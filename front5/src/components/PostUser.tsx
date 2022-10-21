import React from "react";
import {Edges, Post, User} from "../store/types";
import {Link} from "./Link";
import styles from "./PostUser.module.css";
import {UserAvatar} from "./UserAvatar";
import {useQuery} from "@tanstack/react-query";
import {fetcher} from "../store/fetcher";

export function PostUser(props: { postId: string }) {
    const {data: post} = useQuery<Post>(["/posts/" + props.postId], fetcher);

    const {data: user} = useQuery<User>(["/users/" + post?.userId], fetcher, {
        enabled: !!post,
    });
    const [isVisible, setIsVisible] = React.useState(false);

    const {data: userPosts} = useQuery<Edges>(["/users/" + post?.userId + "/posts"], fetcher, {
        enabled: isVisible,
    });

    let className = styles.userNamePopup;
    if (!isVisible) {
        className += " " + styles.userNamePopup__hidden;
    }


    let userDetails = '...LOADING...';
    if (userPosts && user) {
        userDetails = "Name: " + user.name + ", posts: " + userPosts.totalCount;
    }

    const date = new Date(post?.date || "");
    const dateStr = date.toLocaleString();

    return (
        <div className={styles.userName}
             onMouseEnter={() => setIsVisible(true)}
             onMouseLeave={() => setIsVisible(false)}
        >
            <div className={className}>{userDetails}</div>

            <UserAvatar width={50} url={user?.avatar}/>
            <div className={styles.rightContainer}>
                <Link href={"/users/" + user?.id} className={styles.name}>{user?.name}</Link>
                <Link href={"/posts/" + post?.id} className={styles.href}>{dateStr}</Link>
            </div>
        </div>
    )
}
