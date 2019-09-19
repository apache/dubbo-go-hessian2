// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package test;

public class TestString {

    public static String getEmojiTestString() {
        // 0x0001_f923
        String s = "emoji\uD83E\uDD23";

        // see: http://www.unicode.org/glossary/#code_point
        int[] ucs4 = new int[]{0x0010_ffff};
        String maxUnicode = new String(ucs4, 0, ucs4.length);

        return s + ",max" + maxUnicode;
    }
}
