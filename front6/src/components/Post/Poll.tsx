import * as types from "../../api/api";
import React from "react";
import * as styles from "./Poll.module.css";
import {useResources} from "../../store/resources";

export function Poll(props: {
    poll: types.Poll,
    postId: string,
}) {
    const resources = useResources();
    const poll = props.poll;

    let isVoted = false;
    for (let answer of poll.answers) {
        isVoted = isVoted || answer.isVoted;
    }

    const onVote = (answerId: string) => {
        if (isVoted) return;

        types.pollsVote({
            pollId: poll.id,
            answerIds: [answerId],
        }).then(() => {
            let post = structuredClone(resources.posts[props.postId]) as types.Post;

            if (!post.poll) {
                return;
            }

            for (let answer of post.poll.answers) {
                if (answer.id == answerId) {
                    answer.isVoted = true;
                    answer.voted = (answer.voted || 0) + 1;
                }
            }

            resources.setPost(post);
        })
    }

    const onDeleteVote = (e: React.MouseEvent<HTMLAnchorElement>) => {
        e.preventDefault();

        types.pollsDeleteVote({
            pollId: poll.id,
        }).then(() => {
            let post = structuredClone(resources.posts[props.postId]) as types.Post;

            if (!post.poll) {
                return;
            }

            for (let answer of post.poll.answers) {
                if (answer.isVoted) {
                    answer.isVoted = false;
                    answer.voted = (answer.voted || 0) - 1;
                }
            }

            resources.setPost(post);
        })
    }


    return <div className={styles.poll}>
        <div className={styles.question}>
            {poll.question}
        </div>

        <div>
            {poll.answers.map(answer => {
                let className = styles.answer;
                if (answer.isVoted) {
                    className += " " + styles.votedAnswer;
                }

                return <div className={className} onClick={() => onVote(answer.id)}>
                    {answer.answer}
                    <div className={styles.votersCount}>{answer.voted || 0}</div>
                </div>
            })}
            {isVoted && <a href="#" className={styles.deleteVote} onClick={onDeleteVote}>Delete vote</a>}
        </div>
    </div>
}
