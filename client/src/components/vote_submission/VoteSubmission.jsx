import React from 'react';
import { useState } from 'react';
import { InputGroup, FormControl, Button } from 'react-bootstrap';
import axios from 'axios';
import './voteSubmission.css';

const VoteSubmission = () => {
    const [dto, setDto] = useState({
        voter_ic_no: "",
        voter_name: "",
        party: ""
    });

    const [errorMsg, setErrorMsg] = useState("");
    const [successMsg, setSuccessMsg] = useState("");

    const handleVoterIcNo = (event) => {
        setDto({ ...dto, voter_ic_no: event.target.value });
    }

    const handleVoterName = (event) => {
        setDto({ ...dto, voter_name: event.target.value });
    }

    const handleParty = (event) => {
        setDto({ ...dto, party: event.target.value });
    }

    const isNotBlank = (input) => {
        if((input?.trim()?.length || 0) > 0) {
            return true
        }
        return false;
    }

    const handleSubmitForm = async (e) => {
        e.preventDefault();
        setErrorMsg("");
        setSuccessMsg("");
        if(!isNotBlank(dto.voter_ic_no) || !isNotBlank(dto.voter_name) || !isNotBlank(dto.party)) {
            setErrorMsg("Please enter every inputs");
            return
        }

        let regex = /^(\d{10}|\d{12})$/;
        let match = regex.test(dto.voter_ic_no);
        if(!match) {
            setErrorMsg("Please enter valid IC number");
            return
        }

        axios.post(`http://localhost:8000/api/create-vote`, dto).then(
            () => {
                setSuccessMsg("Successfully submitted for: " + dto.voter_ic_no + " !!!");
                setDto({
                    voter_ic_no: "",
                    voter_name: "",
                    party: ""
                });
            }
        ).catch(
            (error) => {
                console.log('Error: ', error);
                if(error.response.data.error) {
                    setErrorMsg(error.response.data.error);
                }
            }
        )
    }

    return (
        <div className='submission__form'>
            <div className="submission__form-inputs grid">
                <div className="submission__form-div">
                    <div className="submission__input-text">
                        <InputGroup.Text>
                            IC Number
                        </InputGroup.Text>
                    </div>
                    <div className="submission__input-input">
                        <FormControl 
                            type='text' 
                            value={dto.voter_ic_no} 
                            onChange={(e) => handleVoterIcNo(e)} 
                            placeholder='Please enter your IC No without dash (-)'
                            className='submission__input-formcontrol'
                        />
                    </div>                    
                </div>
                <div className="submission__form-div">
                    <div className="submission__input-text">
                        <InputGroup.Text>
                            Name
                        </InputGroup.Text>  
                    </div>
                    <div className="submission__input-input">
                        <FormControl 
                            type='text' 
                            value={dto.voter_name} 
                            onChange={(e) => handleVoterName(e)} 
                            placeholder='Please enter your name'
                            className='submission__input-formcontrol'
                        />
                    </div>                    
                </div>
                <div className="submission__form-div">
                    <div className="submission__input-text">
                        <InputGroup.Text>
                            Party
                        </InputGroup.Text>
                    </div>
                    <div className="submission__input-input">
                        <select value={dto.party} onChange={(e) => handleParty(e)} placeholder='Parties'>
                            <option disabled={true} value=''>
                                --Choose a party--
                            </option>
                            <option key='PH' value='PH'>Pakatan Harapan</option>
                            <option key='BN' value='BN'>Barisan Nasional</option>  
                            <option key='PN' value='PN'>Perikatan Nasional</option>  
                        </select>
                    </div>                    
                </div>
            </div> 
            <Button onClick={handleSubmitForm} variant='primary' className='button button--flex'>Submit</Button>
            {errorMsg ? (<p className='errormessage'>{errorMsg}</p>) : null}
            {successMsg ? (<p>{successMsg}</p>) : null}
        </div>
    )
}

export default VoteSubmission