package com.example.AgentApp.mapper;


import com.example.AgentApp.dto.RegistrationRequestDto;
import com.example.AgentApp.enums.Gender;
import com.example.AgentApp.enums.UserRole;
import com.example.AgentApp.model.User;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Component;

import java.text.ParseException;
import java.text.SimpleDateFormat;
import java.util.Date;

@Component
public class UserMapper {

    @Autowired
    private PasswordEncoder passwordEncoder;

    public User mapToUser(RegistrationRequestDto dto) throws ParseException {
        User user = new User();
        user.setConfirmed(false);
        user.setFirstName(dto.getFirstName());
        user.setLastName(dto.getLastName());
        user.setDateOfBirth(getDateOfBirthFromRequest(dto.getDateOfBirth()));
        user.setGender(getGenderFromRequest(dto.getGender()));
        user.setPassword(passwordEncoder.encode(dto.getPassword()));
        user.setEmail(dto.getEmail());
        user.setRecoveryEmail(dto.getRecoveryEmail());
        user.setUsername(dto.getUsername());
        user.setPhoneNumber(dto.getPhoneNumber());
        user.setRole(UserRole.REGISTERED_USER);
        return user;
    }

    private Date getDateOfBirthFromRequest(String dateOfBirth) throws ParseException {
        return new SimpleDateFormat("MM/dd/yyyy").parse(dateOfBirth);
    }

    private Gender getGenderFromRequest(String gender) {
        if(gender == Gender.FEMALE.toString()){
            return Gender.FEMALE;
        }else{
            return Gender.MALE;
        }

    }
}