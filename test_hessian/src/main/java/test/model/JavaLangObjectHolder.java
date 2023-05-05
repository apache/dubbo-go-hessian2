/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package test.model;

import java.io.Serializable;

/**
 * @author tiltwind
 */
public class JavaLangObjectHolder implements Serializable {
    private Integer fieldInteger;
    private Long fieldLong;
    private Boolean fieldBoolean;
    private Double fieldDouble;
    private Float fieldFloat;
    private Short fieldShort;
    private Byte fieldByte;
    private Character fieldCharacter;

    public Integer getFieldInteger() {
        return fieldInteger;
    }

    public void setFieldInteger(Integer fieldInteger) {
        this.fieldInteger = fieldInteger;
    }

    public Long getFieldLong() {
        return fieldLong;
    }

    public void setFieldLong(Long fieldLong) {
        this.fieldLong = fieldLong;
    }

    public Boolean getFieldBoolean() {
        return fieldBoolean;
    }

    public void setFieldBoolean(Boolean fieldBoolean) {
        this.fieldBoolean = fieldBoolean;
    }

    public Double getFieldDouble() {
        return fieldDouble;
    }

    public void setFieldDouble(Double fieldDouble) {
        this.fieldDouble = fieldDouble;
    }

    public Float getFieldFloat() {
        return fieldFloat;
    }

    public void setFieldFloat(Float fieldFloat) {
        this.fieldFloat = fieldFloat;
    }

    public Short getFieldShort() {
        return fieldShort;
    }

    public void setFieldShort(Short fieldShort) {
        this.fieldShort = fieldShort;
    }

    public Byte getFieldByte() {
        return fieldByte;
    }

    public void setFieldByte(Byte fieldByte) {
        this.fieldByte = fieldByte;
    }

    public Character getFieldCharacter() {
        return fieldCharacter;
    }

    public void setFieldCharacter(Character fieldCharacter) {
        this.fieldCharacter = fieldCharacter;
    }
}
