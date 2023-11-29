import {Poll as ApiPoll, pollsDeleteVote, pollsVote} from "../../api/api";
import React from "react";
import * as styles from "./Poll.module.css";

export function Poll(props: {
    poll: ApiPoll
}) {
    const poll = props.poll;

    let isVoted = false;
    for (let answer of poll.answers) {
        isVoted = isVoted || answer.isVoted;
    }

    const onVote = (answerId: string) => {
        if (isVoted) return;

        pollsVote({
            pollId: poll.id,
            answerIds: [answerId],
        })
    }

    const onDeleteVote = (e: React.MouseEvent<HTMLAnchorElement>) => {
        e.preventDefault();

        pollsDeleteVote({
            pollId: poll.id,
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
