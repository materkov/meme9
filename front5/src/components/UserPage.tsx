import React, {useEffect} from "react";
import * as types from "../api/types";
import {uploadApi} from "../api/api";
import {ComponentPost} from "./Post";
import styles from "./UserPage.module.css";
import {localizeCounter} from "../utils/localize";
import {UserAvatar} from "./UserAvatar";
import {Global} from "../store/store";
import {connect} from "react-redux";
import {edit, follow, loadUserPage, unfollow, usersSetAvatar} from "../store/actions/users";

interface Props {
    user: types.User;

    posts: number;
    postIds: string[];
    followers: number;
    following: number;
    isFollowing: boolean;

    viewerId: string;
}

function Component(props: Props) {
    const [isSuccess, setIsSuccess] = React.useState(false);
    const userId = location.pathname.substring(7);
    useEffect(() => {
        loadUserPage(userId).then(() => setIsSuccess(true));
    }, []);

    const [userName, setUserName] = React.useState("");
    const [userNameUpdated, setUserNameUpdated] = React.useState(false);

    const [avatarUploading, setAvatarUploading] = React.useState(false);

    const editName = () => {
        edit({
            userId: userId,
            name: userName,
        }).then(() => {
            setUserNameUpdated(true);
        });
    }

    const onFollow = () => {
        follow(userId);
    }

    const onUnfollow = () => {
        unfollow(userId);
    }

    const onAvatarUpload = (e: React.ChangeEvent<HTMLInputElement>) => {
        if (!e.target.files) {
            return;
        }

        setAvatarUploading(true);
        const file = e.target.files[0];

        uploadApi(file).then(uploadToken => {
            usersSetAvatar(uploadToken).then(() => {
                setAvatarUploading(false);
            });
        })
    }

    if (!isSuccess) {
        return <>Загрузка ...</>;
    }

    const user = props.user;

    return (
        <div>
            <div className={styles.topBlock}>
                <UserAvatar width={100} userId={userId}/>
                <div className={styles.infoBlock}>
                    <div className={styles.userName}>{user.name}</div>
                    <div className={styles.userBio}>
                        {user.bio}
                    </div>
                    <div className={styles.userCounters}>
                        <div className={styles.userCounter}>
                            <b>{props.posts}</b> {localizeCounter(props.posts, "публикация", "публикации", "публикаций")}
                        </div>
                        <div className={styles.userCounter}>
                            <b>{props.followers}</b> {localizeCounter(props.followers, "подписчик", "подписчика", "подписчиков")}
                        </div>
                        <div className={styles.userCounter}>
                            <b>{props.following}</b> {localizeCounter(props.following, "подписка", "подписки", "подписок")}
                        </div>
                        <div className={styles.buttonsBlock}>
                            {props.viewerId && props.viewerId != user.id &&
                                <>
                                    {props.isFollowing ?
                                        <button onClick={onUnfollow}>Отписаться</button> :
                                        <button onClick={onFollow}>Подписаться</button>
                                    }
                                </>
                            }
                        </div>
                    </div>
                </div>
            </div>

            {user.id === props.viewerId && <>
                Имя: <input type="text" value={userName} onChange={e => setUserName(e.target.value)}/>
                <button onClick={editName}>Обновить</button>
                {userNameUpdated && <div>Имя успешно обновлено</div>}

                <br/><br/>
                Поменять аватарку:
                <br/>
                {!avatarUploading && <input type={"file"} onChange={onAvatarUpload}/>}
                {avatarUploading && <span>Загружаем аватар...</span>}
            </>}

            {props.postIds.map(postId => (
                <ComponentPost key={postId} id={postId}/>
            ))}

            {/*
            {hasNextPage && <button onClick={() => fetchNextPage()}>Показать еще</button>}
            */}
        </div>
    )
}

export const UserPage = connect((state: Global) => {
    const userId = location.pathname.substring(7);

    return {
        user: state.users.byId[userId],
        posts: state.users.postsCount[userId],
        followers: state.users.followersCount[userId],
        following: state.users.followingCount[userId],
        isFollowing: state.users.isFollowing[userId],
        viewerId: state.viewer.id,
        postIds: state.users.posts[userId],
    } as Props
})(Component);
