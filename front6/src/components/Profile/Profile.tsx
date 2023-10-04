import React, {useEffect} from "react";
import {useProfile} from "../../store/profile";
import {Post} from "../Post/Post";
import * as styles from "./Profile.module.css";
import {SubscribeAction, User, usersFollow, usersSetStatus} from "../../api/api";
import {useGlobals} from "../../store/globals";
import {useResources} from "../../store/resources";

export function Profile() {
    const userId = document.location.pathname.substring(7);
    const profileState = useProfile();
    const globals = useGlobals();
    const resources = useResources();

    const [status, setStatus] = React.useState("");

    useEffect(() => {
        profileState.fetch(userId);
    }, []);

    const updateStatus = () => {
        usersSetStatus({status: status})
            .then(() => {
                const user: User = structuredClone(resources.users[userId]);
                user.status = status;
                resources.setUser(user);
            });
    };

    const follow = () => {
        usersFollow({
            targetId: userId,
            action: user.isFollowing ? SubscribeAction.UNFOLLOW : SubscribeAction.FOLLOW,
        }).then(() => {
            const user: User = structuredClone(resources.users[userId]);
            user.isFollowing = !user.isFollowing;
            resources.setUser(user);
        });
    }

    if (!resources.users[userId]) {
        return <div>Loading....</div>
    }

    const user = resources.users[userId];

    const postIds = profileState.postIds[userId] || [];
    const posts = postIds.map(postId => resources.posts[postId]);

    return <div>
        <h1 className={styles.userName}>{user.name}</h1>
        <div>{user.status}</div>

        {globals.viewerId === user.id && <>
        <textarea placeholder="Your text status..." className={styles.statusInput} value={status}
                  onChange={e => setStatus(e.target.value)}></textarea>
            <button onClick={updateStatus}>Update</button>
        </>
        }

        {globals.viewerId && globals.viewerId !== user.id && <>
            {user.isFollowing ?
                <button onClick={follow}>Unfollow user</button> :
                <button onClick={follow}>Follow user</button>
            }

        </>}

        <hr/>

        {posts.map(post => (
            <Post post={post} key={post.id}/>
        ))}
    </div>
}


