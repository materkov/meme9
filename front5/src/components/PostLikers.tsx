import React, {useEffect} from "react";
import styles from "./PostLikers.module.css";
import {Link} from "./Link";
import {UserAvatar} from "./UserAvatar";
import {Global} from "../store/store";
import {connect} from "react-redux";
import * as types from "../api/types";
import {loadLikers} from "../store/actions/posts";

interface Props {
    postId: string;
    count: number;
    isLiked: boolean;
    likers: types.User[];
}

function Component(props: Props) {
    const [isLoading, setIsLoading] = React.useState(true);

    useEffect(() => {
        loadLikers(props.postId).then(() => setIsLoading(false));
    }, []);

    return <div className={styles.list}>
        {isLoading && <>Загрузка...</>}
        {!isLoading && !props.count && <>Никто не полайкал.</>}

        {!isLoading && (props.likers).map(user => (
            <Link className={styles.item} href={"/users/" + user.id} key={user.id}>
                <UserAvatar width={40} userId={user.id}/>
                <div className={styles.name}>{user.name}</div>
            </Link>
        ))}
    </div>
}

export const PostLikers = connect((state: Global, ownProps: { id: string }) => {
    return {
        postId: ownProps.id,
        count: state.posts.likesCount[ownProps.id],
        isLiked: state.posts.isLiked[ownProps.id],
        likers: (state.posts.likers[ownProps.id] || []).map(userId => state.users.byId[userId]),
    } as Props
})(Component);
