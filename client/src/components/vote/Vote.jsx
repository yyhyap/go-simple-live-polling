import React from 'react';
import VoteCount from '../vote_count/VoteCount';
import VoteSubmission from '../vote_submission/VoteSubmission';
import './vote.css';

const Vote = () => {
    return (
        <>
            <VoteCount />
            <VoteSubmission />
        </>
    )
}

export default Vote