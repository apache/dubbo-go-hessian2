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
 * Mirrors the Go struct {@code Int64ThenStringStruct}: a {@code Long} field
 * followed by a {@code String} field. Used to cross-validate that a Go provider
 * encodes a {@code *int64} field as a hessian long (non-nil) or hessian null
 * (nil) without shifting the following field, see apache/dubbo-go#2410.
 */
public class Int64ThenString implements Serializable {
    private Long total;
    private String note;

    public Long getTotal() {
        return total;
    }

    public void setTotal(Long total) {
        this.total = total;
    }

    public String getNote() {
        return note;
    }

    public void setNote(String note) {
        this.note = note;
    }
}
