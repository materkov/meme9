import React from "react";
import {Link} from "../Link/Link";
import * as schema from "../../api/api2";
import {Post} from "../Feed/Post/Post";
import {GlobalStore, GlobalStoreContext} from "../../Context";

export function Profile(props: { data: schema.ProfileRenderer }) {
    const data = props.data;
    const theme = React.useContext(GlobalStoreContext) as GlobalStore;

    const onSubscribe = () => {
        theme.followUser(data.id);
    };

    const onUnsubscribe = () => {
        theme.unfollowUser(data.id);
    };

    return <div>
        Profile page
        <br/>
        <Link href={"/"}>Go to feed</Link>
        <br/>
        <img alt="" src={data.avatar}/>
        <br/>
        {data.isFollowing ? 'Вы подписаны' : 'Вы не подписаны'}
        {data.isFollowing ?
            <button onClick={onUnsubscribe}>Отписаться</button> :
            <button onClick={onSubscribe}>Подписаться</button>
        }
        <br/>
        <h2>{data.name}</h2>
        {data.posts.map(post => (
            <Post key={post.id} data={post}/>
        ))}
    </div>;
}
