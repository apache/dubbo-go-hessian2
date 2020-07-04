package test;

import java.time.*;

/**
 * date 2020/7/2 11:07 <br/>
 * description class <br/>
 * add java8 package java.time Object method
 * 添加 java8 java.time 包下的各种序列化对象的获取 方便测试
 *
 * @author danao
 * @version 1.0
 * @since 1.0
 */
public class TestJava8Time {


    public static Duration java8_Duration() {
        return Duration.ZERO;
    }

    public static Instant java8_Instant() {
        return Instant.ofEpochMilli(100000L);
    }

    public static LocalDate java8_LocalDate() {
        return LocalDate.of(2020, 6, 6);
    }

    public static LocalDateTime java8_LocalDateTime() {
        return LocalDateTime.of(2020, 6, 6, 6, 6);
    }

    public static LocalTime java8_LocalTime() {
        return LocalTime.of(6, 6);
    }

    public static MonthDay java8_MonthDay() {
        return MonthDay.of(6, 6);
    }

    public static OffsetDateTime java8_OffsetDateTime() {
        return OffsetDateTime.of(2020, 6, 6, 6, 6, 6, 6, ZoneOffset.MIN);
    }

    public static OffsetTime java8_OffsetTime() {
        return OffsetTime.of(6, 6, 6, 6, ZoneOffset.MIN);
    }

    public static Period java8_Period() {
        return Period.of(2020, 6, 6);
    }

    public static Year java8_Year() {
        return Year.of(2020);
    }

    public static YearMonth java8_YearMonth() {
        return YearMonth.of(2020, 6);
    }

    public static ZonedDateTime java8_ZonedDateTime() {
        ZonedDateTime of = ZonedDateTime.of(java8_LocalDateTime(), java8_ZoneId());
        return of;
    }

    public static ZoneId java8_ZoneId() {
        return ZoneId.of("1000");
    }


    public static ZoneOffset java8_ZoneOffset() {
        return ZoneOffset.of("1000");
    }

    public static void main(String[] args) {

    }
}
