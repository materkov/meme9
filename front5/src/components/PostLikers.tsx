import React, {useEffect} from "react";
import styles from "./PostLikers.module.css";
import {Link} from "./Link";
import {UserAvatar} from "./UserAvatar";
import {Global} from "../store2/store";
import {actions} from "../store2/actions";
import {connect} from "react-redux";
import * as types from "../store/types";

interface Props {
    postId: string;
    count: number;
    isLiked: boolean;
    likers: types.User[];
}

function Component(props: Props) {
    const [isLoading, setIsLoading] = React.useState(true);

    useEffect(() => {
        actions.loadLikers(props.postId).then(() => setIsLoading(false));
    }, []);
    //const {data, isLoading, isStale} = useQuery<Edges>(["/posts/" + props.id + "/liked?count=10"], fetcher, {
    //})

    /*const userQueries = useQueries<User[]>({
        queries: (data?.items || []).map(userId => {
            return {
                queryKey: ["/users/" + userId],
                queryFn: fetcher,
            }
        })
    })*/

    //const users: User[] = [];
    //for (let q of userQueries) {
    //    users.push(q.data as User);
    //}

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
