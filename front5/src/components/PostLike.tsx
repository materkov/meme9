import styles from "./PostLike.module.css";
import {Heart} from "./icons/Heart";
import React from "react";
import {HeartRed} from "./icons/HeartRed";
import classNames from "classnames";
import {PostLikers} from "./PostLikers";
import {Global} from "../store2/store";
import {connect} from "react-redux";
import {like, unlike} from "../store2/actions/posts";

type Props = {
    postId: string;
    count: number;
    isLiked: boolean;
    viewerId: string;
}

const Component = (props: Props) => {
    const [likersVisible, setLikersVisible] = React.useState(false);

    const onClick = () => {
        if (!props.viewerId) return;

        if (!props.isLiked) {
            like(props.postId);
        } else {
            unlike(props.postId);
        }
    }

    return <div className={styles.likeBtn} onClick={onClick}
                onMouseEnter={() => setLikersVisible(true)}
                onMouseLeave={() => setLikersVisible(false)}
    >
        {props.isLiked ?
            <HeartRed className={styles.likeIcon}/> :
            <Heart className={styles.likeIcon}/>
        }

        {props.count > 0 &&
            <div className={styles.likeText}>{props.count}</div>
        }

        {likersVisible &&
            <div className={classNames({
                [styles.likersPopup]: true,
                [styles.likersPopup__visible]: true,
            })}>
                <PostLikers id={props.postId}/>
            </div>
        }
    </div>
}


export const PostLike = connect((state: Global, ownProps: { id: string }) => {
    return {
        postId: ownProps.id,
        count: state.posts.likesCount[ownProps.id],
        isLiked: state.posts.isLiked[ownProps.id],
        viewerId: state.viewer.id,
    } as Props
})(Component);