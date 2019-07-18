package test.model;

import java.io.Serializable;
import java.util.Date;

public class DateDemo implements Serializable {
    private String name;
    private Date date;

    public String getName() {
        return name;
    }

    public DateDemo() {}

    public DateDemo(String name,Date date) {
        this.name = name;
        this.date = date;
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

}