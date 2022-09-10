import React, {useEffect} from "react";
import {api, Post, User} from "../store/types";
import {ComponentPost} from "./Post";

export function UserPage() {
    const [user, setUser] = React.useState<User>();
    const [viewerId, setViewerId] = React.useState("");
    const [posts, setPosts] = React.useState<Post[]>([]);
    const [loaded, setLoaded] = React.useState(false);

    const [userName, setUserName] = React.useState("");
    const [userNameUpdated, setUserNameUpdated] = React.useState(false);

    useEffect(() => {
        const f = new FormData();
        f.set("id", location.pathname.substring(7));

        api("/userPage", {
            id: location.pathname.substring(7)
        }).then(r => {
            const [user, viewerId] = r;
            setUser(user);
            setViewerId(viewerId);
            setUserName(user.name);

            // TODO strange
            for (let post of user.posts) {
                post.user = user;
            }

            setPosts(user.posts);
            setLoaded(true);
        })
    }, []);

    const editName = () => {
        api("/userEdit", {
            id: viewerId,
            name: userName,
        }).then(() => setUserNameUpdated(true));
    }

    if (!loaded || !user) {
        return <>Загрузка ...</>;
    }

    return (
        <div>
            {user.name}
            <hr/>
            {user.id === viewerId && <>
                Имя: <input type="text" value={userName} onChange={e => setUserName(e.target.value)}/>
                <button onClick={editName}>Обновить</button>
                {userNameUpdated && <div>Имя успешно обновлено</div>}
            </>}
            {posts.map(post => <ComponentPost key={post.id} post={post}/>)}
        </div>
    )
}
