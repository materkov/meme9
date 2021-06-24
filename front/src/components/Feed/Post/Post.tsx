import React from "react";
import styles from "./Post.module.css";
import * as schema from "../../../api/api2";
import {Link} from "../../Link/Link";
import {Heart} from "../../../icons/Heart";
import {HeartRed} from "../../../icons/HeartRed";
import {GlobalStoreContext} from "../../../Context";
import {Comment as CommentIcon} from "../../../icons/Comment";
import {Comment} from "../Comment/Comment";

export interface Props {
    data: schema.Post;
}

interface State {
    likeHovered: boolean;
}

export class Post extends React.Component<Props, State> {
    static contextType = GlobalStoreContext;

    state: State = {
        likeHovered: false,
    };

    onLikeHover = () => {
        this.setState({likeHovered: true});
    }

    onLikeBlur = () => {
        this.setState({likeHovered: false});
    }

    onToggleLike = () => {
        if (!this.props.data.canLike) {
            alert('Нужно авторизоваться, чтобы полайкать');
            return;
        }

        this.context.togglePostLike(this.props.data.id);
    }

    render() {
        const data = this.props.data;

        return (
            <div className={styles.Post}>
                <div className={styles.Header}>
                    <img className={styles.AuthorAvatar} alt="" src={data.authorAvatar}/>
                    <div>
                        <Link className={styles.Author} href={data.authorUrl}>{data.authorName}</Link>
                        <Link href={data.url} className={styles.Date}>{data.dateDisplay}</Link>
                    </div>
                </div>

                <div className={styles.Text}>{data.text}</div>

                {data.imageUrl &&
                <img className={styles.Image} alt="" src={data.imageUrl}/>
                }

                <div className={styles.LikeContainer} onMouseEnter={this.onLikeHover} onMouseLeave={this.onLikeBlur}
                     onClick={this.onToggleLike}
                >
                    {(this.state.likeHovered || data.isLiked) ?
                        <HeartRed className={styles.Like}/> :
                        <Heart className={styles.Like}/>
                    }
                    <div className={styles.LikeCounter}>{data.likesCount}</div>
                </div>

                <Link href={data.url} className={styles.CommentContainer}>
                    <CommentIcon className={styles.CommentIcon}/>
                    <div className={styles.CommentCounter}>{data.commentsCount}</div>
                </Link>

                <TopComment post={data}/>
            </div>
        );
    }

}

const TopComment = (props: { post: schema.Post }) => {
    if (!props.post.topComment) {
        return null;
    }

    return <div className={styles.TopCommentContainer}><Comment data={props.post.topComment}/></div>
}
