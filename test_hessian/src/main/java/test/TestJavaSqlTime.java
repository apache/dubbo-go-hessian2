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

package test;

import java.sql.Date;
import java.sql.Time;

public class TestJavaSqlTime {

    public Object javaSql_decode_date() {
        Date date = Date.valueOf("2020-08-09");
        return date;
    }

    public Object javaSql_decode_time() {
        Time time = new Time(852095746000L);
        return time;
    }


    public TestJavaSqlTime() {
    }

    public static Object javaSql_encode_time(Object v) {
        return v.equals(new Time(852095746000L));
    }

    public boolean javaSql_encode_date(Object v) {
        Date date = Date.valueOf("2020-08-09");
        return date.equals(v);
    }

}
