package models

import (
    "go.mongodb.org/mongo-driver/bson/primitive"
    // "time"
)

// type Sim struct {
//     ID                   primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
//     Name                 string             `bson:"name" json:"name"`
//     Number               string             `bson:"number" json:"number"`
//     LastRechargeDate     string             `bson:"last_recharge_date" json:"last_recharge_date"`
//     RechargeValidity     string             `bson:"recharge_validity" json:"recharge_validity"`
//     IncomingCallValidity string             `bson:"incoming_call_validity" json:"incoming_call_validity"`
//     SimExpiry            string             `bson:"sim_expiry" json:"sim_expiry"`
// }

type Sim struct {
    ID                   primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
    Name                 string             `bson:"name" json:"name"`
    Number               string             `bson:"number" json:"number"`
    LastRechargeDate     string             `bson:"last_recharge_date" json:"last_recharge_date"`
    RechargeValidity     string             `bson:"recharge_validity" json:"recharge_validity"`
    IncomingCallValidity string             `bson:"incoming_call_validity" json:"incoming_call_validity"`
    SimExpiry            string             `bson:"sim_expiry" json:"sim_expiry"`
}