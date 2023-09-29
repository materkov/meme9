import React, {useEffect} from "react";
import {useProfile} from "../../store/profile";
import {Post} from "../Post/Post";
import * as styles from "./Profile.module.css";
import {usersSetStatus} from "../../api/api";
import {useGlobals} from "../../store/globals";

export function Profile() {
    const userId = document.location.pathname.substring(7);
    const profileState = useProfile();
    const globals = useGlobals();

    const [status, setStatus] = React.useState("");

    useEffect(() => {
        profileState.fetch(userId);
    }, []);

    const updateStatus = () => {
        usersSetStatus({status: status})
            .then(() => profileState.setStatus(globals.viewerId, status));
    };

    if (!profileState.users[userId]) {
        return <div>Loading....</div>
    }

    const user = profileState.users[userId];

    return <div>
        <h1 className={styles.userName}>{user.name}</h1>
        <div>{user.status}</div>

        {globals.viewerId === user.id && <>
        <textarea placeholder="Your text status..." className={styles.statusInput} value={status}
                  onChange={e => setStatus(e.target.value)}></textarea>
            <button onClick={updateStatus}>Update</button>
        </>
        }

        <hr/>

        {(profileState.posts[userId] || []).map(post => (
            <Post post={post} key={post.id}/>
        ))}
    </div>
}


