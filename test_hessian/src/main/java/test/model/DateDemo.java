package test.model;

import java.io.Serializable;
import java.util.Date;

public class DateDemo implements Serializable {
    private String name;
    private Date date;
    private Date date1;

    public String getName() {
        return name;
    }

    public DateDemo() {}

    public DateDemo(String name,Date date,Date date1) {
        this.name = name;
        this.date = date;
        this.date1 = date1;
    }

    public void setName(String name) {
        this.name = name;
    }
    public Date getDate() {
        return date;
    }

    public void setDate(Date date) {
        this.date = date;
    }
    public Date getDate1() {
        return date1;
    }

    public void setDate1(Date date1) {
        this.date1 = date1;
    }

}