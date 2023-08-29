import React, {useEffect} from "react";
import {useProfile} from "../../store/profile";
import {Post} from "../Post/Post";

export function Profile() {
    const userId = document.location.pathname.substring(7);
    const profileState = useProfile();

    useEffect(() => {
        profileState.fetch(userId);
    }, []);

    if (!profileState.user.id) {
        return <div>Loading....</div>
    }

    return <div>
        <h1>{profileState.user.name}</h1>
        <hr/>

        {profileState.posts.map(post => (
            <Post post={post} key={post.id}/>
        ))}
    </div>
}


