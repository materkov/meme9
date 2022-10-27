import React from "react";
import {api, Edges, User, UserPostsConnection, Viewer} from "../store/types";
import {ComponentPost} from "./Post";
import styles from "./UserPage.module.css";
import {localizeCounter} from "../utils/localize";
import {UserAvatar} from "./UserAvatar";
import {useQuery} from "@tanstack/react-query";
import {fetcher, queryClient} from "../store/fetcher";

export function UserPage() {
    const userId = location.pathname.substring(7);
    const {data: user, isSuccess} = useQuery<User>(["/users/" + userId], fetcher);
    const {data: posts} = useQuery<Edges>(["/users/" + userId + "/posts"], fetcher);

    const followersQueryKey = ["/users/" + userId + "/followers"];
    const {data: followers} = useQuery<Edges & {isFollowing: boolean}>(followersQueryKey, fetcher);
    const {data: following} = useQuery<Edges>(["/users/" + userId + "/following"], fetcher);
    const {data: viewer} = useQuery<Viewer>(["/viewer"], fetcher);

    const [userName, setUserName] = React.useState("");
    const [userNameUpdated, setUserNameUpdated] = React.useState(false);

    const [avatarUploading, setAvatarUploading] = React.useState(false);

    const editName = () => {
        api("/userEdit", {
            id: userId,
            name: userName,
        }).then(() => {
            setUserNameUpdated(true);

            const queryKey = ["/users/" + userId];
            const user = queryClient.getQueryData<User>(queryKey);
            if (!user) {
                return;
            }
            queryClient.setQueryData(queryKey, {...user, name: userName});
        });
    }

    const onShowMore = () => {
        // TODO
        api("/userPage/posts", {
            id: location.pathname.substring(7),
            cursor: "",
        }).then((result: [UserPostsConnection]) => {
            let r = result[0];
            //setPostsCursor(r.nextCursor || "");
        })
    }

    const onFollow = () => {
        const oldData = queryClient.getQueryData(followersQueryKey);
        if (oldData) {
            queryClient.setQueryData(followersQueryKey, {
                ...oldData,
                totalCount: (oldData.totalCount || 0) + 1,
                isFollowing: true,
            });
        }

        api("/userFollow", {
            id: userId,
        }).then(() => {
            queryClient.invalidateQueries(followersQueryKey);
        })
    }

    const onUnfollow = () => {
        const oldData = queryClient.getQueryData(followersQueryKey);
        if (oldData) {
            queryClient.setQueryData(followersQueryKey, {
                ...oldData,
                totalCount: (oldData.totalCount || 0) - 1,
                isFollowing: false,
            });
        }

        api("/userUnfollow", {
            id: userId,
        }).then(() => {
            queryClient.invalidateQueries(followersQueryKey);
        })
    }

    const onAvatarUpload = (e: React.ChangeEvent<HTMLInputElement>) => {
        if (!e.target.files) {
            return;
        }

        setAvatarUploading(true);
        const file = e.target.files[0];

        api("/uploadAvatar", {
            file: file,
        }).then((resp) => {
            //setUser({...user, avatar: resp.avatar});
            queryClient.invalidateQueries(["/users/" + userId]);
            setAvatarUploading(false);
        })

    }

    if (!isSuccess) {
        return <>Загрузка ...</>;
    }

    return (
        <div>
            <div className={styles.topBlock}>
                <UserAvatar width={100} url={user.avatar}/>
                <div className={styles.infoBlock}>
                    <div className={styles.userName}>{user.name}</div>
                    <div className={styles.userBio}>
                        {user.bio}
                    </div>
                    <div className={styles.userCounters}>
                        <div className={styles.userCounter}>
                            <b>{posts?.totalCount}</b> {localizeCounter(posts?.totalCount || 0, "публикация", "публикации", "публикаций")}
                        </div>
                        <div className={styles.userCounter}>
                            <b>{followers?.totalCount || 0}</b> {localizeCounter(following?.totalCount || 0, "подписчик", "подписчика", "подписчиков")}
                        </div>
                        <div className={styles.userCounter}>
                            <b>{following?.totalCount || 0}</b> {localizeCounter(following?.totalCount || 0, "подписка", "подписки", "подписок")}
                        </div>
                        <div className={styles.buttonsBlock}>
                            {viewer?.viewerId && viewer.viewerId != user.id &&
                                <>
                                    {followers?.isFollowing ?
                                        <button onClick={onUnfollow}>Отписаться</button> :
                                        <button onClick={onFollow}>Подписаться</button>
                                    }
                                </>
                            }
                        </div>
                    </div>
                </div>
            </div>

            {user.id === viewer?.viewerId && <>
                Имя: <input type="text" value={userName} onChange={e => setUserName(e.target.value)}/>
                <button onClick={editName}>Обновить</button>
                {userNameUpdated && <div>Имя успешно обновлено</div>}

                <br/><br/>
                Поменять аватарку:
                <br/>
                {!avatarUploading && <input type={"file"} onChange={onAvatarUpload}/>}
                {avatarUploading && <span>Загружаем аватар...</span>}
            </>}

            {posts?.items?.map(postId => (
                <ComponentPost key={postId} id={postId}/>
            ))}

            {<button onClick={onShowMore}>Показать еще</button>}
        </div>
    )
}
