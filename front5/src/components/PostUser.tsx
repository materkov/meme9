import React from "react";
import {Link} from "./Link";
import styles from "./PostUser.module.css";
import {UserAvatar} from "./UserAvatar";
import {Global} from "../store2/store";
import * as types from "../store/types";
import {connect} from "react-redux";

interface Props {
    post: types.Post;
    user: types.User;
}

function Component(props: Props) {
    const [isVisible, setIsVisible] = React.useState(false);

    //const {data: userPosts} = useQuery<Edges>(["/users/" + post?.userId + "/posts"], fetcher, {
    //    enabled: isVisible,
    //});
    const userPosts = null;

    let className = styles.userNamePopup;
    if (!isVisible) {
        className += " " + styles.userNamePopup__hidden;
    }


    let userDetails = '...LOADING...';
    if (userPosts && props.user) {
        userDetails = "Name: " + props.user.name + ", posts: " + (userPosts || 0);
    }

    const date = new Date(props.post.date || "");
    const dateStr = date.toLocaleString();

    return (
        <div className={styles.userName}
             onMouseEnter={() => setIsVisible(true)}
             onMouseLeave={() => setIsVisible(false)}
        >
            <div className={className}>{userDetails}</div>

            <UserAvatar width={50} userId={props.post.userId || ""}/>
            <div className={styles.rightContainer}>
                <Link href={"/users/" + props.user.id} className={styles.name}>{props.user.name}</Link>
                <Link href={"/posts/" + props.post.id} className={styles.href}>{dateStr}</Link>
            </div>
        </div>
    )
}

export const PostUser = connect((state: Global, ownProps: { postId: string }) => {
    return {
        post: state.posts.byId[ownProps.postId],
        user: state.users.byId[state.posts.byId[ownProps.postId].userId],
    }
})(Component);
